import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from '@/pages/Home';
import About from '@/pages/About';
import Header from "@/components/header"
import Footer from '@/components/footer';
import SubscriptionPopup from '@/components/SubscriptionPopup';
import './App.css';

interface Data {
  message: string;
}

const App: React.FC = () => {
  const [data, setData] = useState<Data | null>(null);
  const [showPopup, setShowPopup] = useState(true); // Manage visibility of the pop-up

  useEffect(() => {
    fetch('https://sellcustombackend.onrender.com/api/data') // Update with your backend URL
      .then(response => response.json())
      .then(data => setData(data));
  }, []);

  const handleClosePopup = () => {
    setShowPopup(false); // Close the pop-up when user clicks 'No, Thanks'
  };

  return (
    <Router>
      <div className="flex flex-col min-h-screen w-full bg-gray-100 dark:bg-slate-800">
        <Header />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/about" element={<About />} />
          <Route path="/product" element={<About />} />
        </Routes>
        <Footer />
        {data ? <div>{data.message}</div> : 'loading'}
        {showPopup && <SubscriptionPopup onClose={handleClosePopup} />} {/* Display pop-up */}
      </div>
    </Router>
  );
};

export default App;
