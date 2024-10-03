import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from '@/pages/Home';
import AddProduct from './pages/AddProduct';
import AddCategory from './pages/AddCategory';
import About from '@/pages/About';
import Header from "@/components/header"
import Footer from '@/components/footer';
import SubscriptionPopup from '@/components/SubscriptionPopup';
import Cookies from 'js-cookie';
import './App.css';
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';

interface Data {
  message: string;
}
const client = new ApolloClient({
  uri: 'https://sellcustombackend.onrender.com/graphql',
  cache: new InMemoryCache(),  
});

const App: React.FC = () => {
  const [data, setData] = useState<Data | null>(null);
  const [showPopup, setShowPopup] = useState(false);
  const [searchTerm, setSearchTerm] = useState<string>('');

  useEffect(() => {
    fetch('https://sellcustombackend.onrender.com/api/data') // Update with your backend URL
      .then(response => response.json())
      .then(data => setData(data));

    // Check if it's the user's first time visiting the website
    const isFirstVisit = !Cookies.get('firstVisit');
    if (isFirstVisit) {
      // Set a cookie to indicate that the user has visited the website
      Cookies.set('firstVisit', 'true', { expires: 30 }); // Expire in 1 year
      setShowPopup(true);
    }
  }, []);

  const handleSearch = (term: string) => {
    setSearchTerm(term);
  };

  const handleClosePopup = () => {
    setShowPopup(false); // Close the pop-up when user clicks 'No, Thanks'
  };

  return (
    <ApolloProvider client={client}>
      <Router>
        <div className="App">
          <Header onSearch={handleSearch}/>
          <Routes>
            <Route path="/" element={<Home searchTerm={searchTerm}/>} />
            <Route path="/about" element={<About />} />
          </Routes>
          <Footer />
          {data ? (
            <p>{data.message}</p>
          ) : (
            'loading'
          )}
          {showPopup && <SubscriptionPopup onClose={handleClosePopup} />}
          <AddCategory />
          <AddProduct />
        </div>
      </Router>
    </ApolloProvider>
  );
};

export default App;