import React, { useState } from 'react';

interface SubscriptionPopupProps {
  onClose: () => void;
}

const SubscriptionPopup: React.FC<SubscriptionPopupProps> = ({ onClose }) => {
  const [email, setEmail] = useState('');
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Here you would typically handle the email submission, e.g., sending it to your backend
    setSubmitted(true);
    setTimeout(() => onClose(), 2000); // Close the popup after a short delay
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50">
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-full max-w-md">
        <h2 className="text-2xl font-semibold mb-4">Subscribe to Our Newsletter!</h2>
        {submitted ? (
          <p className="text-green-600">Thank you for subscribing!</p>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <p className="text-gray-700 dark:text-gray-300">Enter your email to receive the latest updates and offers.</p>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Your email address"
              className="w-full px-4 py-2 border-b border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-0 focus:border-black dark:focus:border-blue-400"
              required
            />
            <button
              type="submit"
              className="w-full px-4 py-2 bg-black text-white border border-black rounded-lg hover:bg-gray-800 dark:bg-black dark:hover:bg-gray-900"
            >
              Subscribe Now
            </button>
            <button
              type="button"
              onClick={onClose}
              className="w-full px-4 py-2 bg-white border border-black text-black rounded-lg hover:bg-gray-200 dark:bg-gray-600 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              No, Thanks
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

export default SubscriptionPopup;
