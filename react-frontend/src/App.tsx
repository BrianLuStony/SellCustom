import React, { useEffect, useState } from 'react';
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'


interface Data {
  message: string;
}

function App() {
  const [count, setCount] = useState(0);
  const [data, setData] = useState<Data | null>(null);

  useEffect(() => {
    fetch('https://your-backend-url.onrender.com/api/data') // Update with your backend URL
      .then(response => response.json())
      .then(data => setData(data));
  }, []);

  return (
    <>
      <div>
      <h1>Frontend React App</h1>
        {data ? <p>{data.message}</p> : <p>Loading...</p>}
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
