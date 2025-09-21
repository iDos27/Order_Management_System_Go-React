import { useState, useEffect, useRef } from 'react';

const useWebSocket = (url) => {
  const [socket, setSocket] = useState(null);
  const [lastMessage, setLastMessage] = useState(null);
  const [connectionStatus, setConnectionStatus] = useState('Connecting');
  const ws = useRef(null);

  useEffect(() => {
    // Tworzymy połączenie WebSocket
    ws.current = new WebSocket(url);
    
    ws.current.onopen = () => {
      console.log('WebSocket połączony');
      setConnectionStatus('Connected');
      setSocket(ws.current);
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Otrzymano wiadomość WebSocket:', data);
      setLastMessage(data);
    };

    ws.current.onclose = () => {
      console.log('WebSocket rozłączony');
      setConnectionStatus('Disconnected');
    };

    ws.current.onerror = (error) => {
      console.error('Błąd WebSocket:', error);
      setConnectionStatus('Error');
    };

    // Cleanup
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [url]);

  return { socket, lastMessage, connectionStatus };
};

export default useWebSocket;