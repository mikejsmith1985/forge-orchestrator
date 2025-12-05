import { useEffect, useRef, useState, useCallback } from 'react';

interface WebSocketMessage {
  type: string;
  payload: any;
}

interface UseWebSocketReturn {
  isConnected: boolean;
  lastMessage: WebSocketMessage | null;
  sendMessage: (message: string | object) => void;
  disconnectedDuration: number;
}

interface UseWebSocketOptions {
  pollingFlowId?: number;
  pollingInterval?: number;
}

const MAX_RECONNECT_ATTEMPTS = 5;
const BASE_RECONNECT_DELAY = 1000; // 1 second
const POLLING_THRESHOLD = 5000; // 5 seconds before starting polling
const DEFAULT_POLLING_INTERVAL = 3000; // 3 seconds

export function useWebSocket(url: string, options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const { pollingFlowId, pollingInterval = DEFAULT_POLLING_INTERVAL } = options;
  
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
  const [disconnectedDuration, setDisconnectedDuration] = useState(0);
  
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<number | undefined>(undefined);
  const disconnectedAtRef = useRef<number | null>(null);
  const disconnectedIntervalRef = useRef<number | undefined>(undefined);
  const pollingIntervalRef = useRef<number | undefined>(undefined);

  // Poll for flow status when WebSocket is disconnected
  const pollFlowStatus = useCallback(async () => {
    if (!pollingFlowId) return;
    
    try {
      const response = await fetch(`/api/flows/${pollingFlowId}/status`);
      if (response.ok) {
        const status = await response.json();
        console.log('Polling received flow status:', status);
        setLastMessage({
          type: 'FLOW_STATUS',
          payload: status,
        });
      }
    } catch (error) {
      console.error('Error polling flow status:', error);
    }
  }, [pollingFlowId]);

  // Start polling when disconnected for > 5 seconds
  useEffect(() => {
    if (!isConnected && disconnectedDuration >= POLLING_THRESHOLD && pollingFlowId) {
      console.log(`WebSocket disconnected for ${disconnectedDuration}ms, starting polling`);
      
      // Start polling immediately
      pollFlowStatus();
      
      // Set up polling interval
      pollingIntervalRef.current = window.setInterval(pollFlowStatus, pollingInterval);
    } else if (isConnected && pollingIntervalRef.current) {
      // Stop polling when WebSocket reconnects
      console.log('WebSocket reconnected, stopping polling');
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = undefined;
    }

    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
      }
    };
  }, [isConnected, disconnectedDuration, pollingFlowId, pollingInterval, pollFlowStatus]);

  // Track disconnected duration
  useEffect(() => {
    if (!isConnected) {
      if (disconnectedAtRef.current === null) {
        disconnectedAtRef.current = Date.now();
      }
      
      disconnectedIntervalRef.current = window.setInterval(() => {
        if (disconnectedAtRef.current !== null) {
          setDisconnectedDuration(Date.now() - disconnectedAtRef.current);
        }
      }, 1000);
    } else {
      disconnectedAtRef.current = null;
      setDisconnectedDuration(0);
      if (disconnectedIntervalRef.current) {
        clearInterval(disconnectedIntervalRef.current);
      }
    }

    return () => {
      if (disconnectedIntervalRef.current) {
        clearInterval(disconnectedIntervalRef.current);
      }
    };
  }, [isConnected]);

  const connect = useCallback(() => {
    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}${url}`;
      
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        reconnectAttemptsRef.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('WebSocket message received:', data);
          setLastMessage(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        setIsConnected(false);
        wsRef.current = null;

        // Exponential backoff reconnection
        if (reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
          const delay = BASE_RECONNECT_DELAY * Math.pow(2, reconnectAttemptsRef.current);
          console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current + 1}/${MAX_RECONNECT_ATTEMPTS})`);
          
          reconnectTimeoutRef.current = window.setTimeout(() => {
            reconnectAttemptsRef.current++;
            connect();
          }, delay);
        } else {
          console.error('Max reconnection attempts reached');
        }
      };
    } catch (error) {
      console.error('Error creating WebSocket connection:', error);
    }
  }, [url]);

  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  const sendMessage = useCallback((message: string | object) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      const messageStr = typeof message === 'string' ? message : JSON.stringify(message);
      wsRef.current.send(messageStr);
    } else {
      console.warn('WebSocket is not connected. Cannot send message.');
    }
  }, []);

  return {
    isConnected,
    lastMessage,
    sendMessage,
    disconnectedDuration,
  };
}
