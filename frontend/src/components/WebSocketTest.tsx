import { useWebSocket } from '../hooks/useWebSocket';

export function WebSocketTestComponent() {
  const { isConnected, lastMessage, sendMessage } = useWebSocket('/ws');

  const handleTestMessage = () => {
    sendMessage({
      type: 'TEST',
      payload: {
        message: 'Test from React Component',
        timestamp: new Date().toISOString(),
      },
    });
  };

  return (
    <div className="p-4 bg-gray-800 rounded-lg">
      <h2 className="text-xl font-bold mb-4 text-white">WebSocket Status</h2>
      
      <div className="mb-4">
        <span className="text-gray-300">Connection Status: </span>
        <span className={isConnected ? 'text-green-400' : 'text-red-400'}>
          {isConnected ? '✅ Connected' : '❌ Disconnected'}
        </span>
      </div>

      {lastMessage && (
        <div className="mb-4">
          <h3 className="text-gray-300 mb-2">Last Message:</h3>
          <pre className="bg-gray-900 p-2 rounded text-xs text-gray-300 overflow-auto">
            {JSON.stringify(lastMessage, null, 2)}
          </pre>
        </div>
      )}

      <button
        onClick={handleTestMessage}
        disabled={!isConnected}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-600 disabled:cursor-not-allowed"
      >
        Send Test Message
      </button>
    </div>
  );
}
