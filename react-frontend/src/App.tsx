import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './pages/Home';
import About from './pages/About';
import './App.css'
import { MainNav } from '@/components/main-nav';

interface Data {
  message: string;
}

function App() {
  const [count, setCount] = useState(0);
  const [data, setData] = useState<Data | null>(null);

  useEffect(() => {
    fetch('https://sellcustombackend.onrender.com/api/data') // Update with your backend URL
      .then(response => response.json())
      .then(data => setData(data));
      console.log(data);
  }, []);

  return (
    <>
     <Router>
        <div className="flex flex-col min-h-screen w-full bg-gray-100 dark:bg-slate-800">
            <MainNav />
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/about" element={<About />} />
              <Route path="/product" element={<About />} />
            </Routes>
            {data ? <div>{data.message}</div> : 'loading'}
          </div>
      </Router>
    </>
  )
}

export default App
