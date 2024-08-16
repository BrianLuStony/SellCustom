import React, { useState } from 'react';
import CustomLink from './custom-link';

const Footer: React.FC = () => {
  const [email, setEmail] = useState('');
  const [submitted, setSubmitted] = useState(false);

  const handleSubscription = (e: React.FormEvent) => {
    e.preventDefault();
    // Handle email subscription logic here
    setSubmitted(true);
    setTimeout(() => setSubmitted(false), 2000); // Reset after a short delay
  };

  return (
    <footer className="w-full dark:bg-slate-800">
      <div className="flex flex-col gap-4 px-4 text-sm sm:flex-row sm:justify-between sm:items-center sm:px-6 sm:mx-auto sm:max-w-3xl sm:h-16">
        <div className="flex flex-col gap-4 sm:flex-row dark:text-gray-300">
          <CustomLink href="https://www.npmjs.com/package/next-auth">NPM</CustomLink>
          <CustomLink href="https://github.com/BrianLuStony/GeniusGrove">Source on GitHub</CustomLink>
          <CustomLink href="/policy">Policy</CustomLink>
        </div>
        <div className="mt-4 sm:mt-0">
          <form onSubmit={handleSubscription} className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Subscribe for updates"
              className="w-full px-4 py-2 border-b border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-0 focus:border-blue-500 dark:focus:border-blue-400"
              required
            />
            <button
              type="submit"
              className="w-full sm:w-auto px-4 py-2 bg-black text-white border border-black rounded-lg hover:bg-gray-800 dark:bg-black dark:hover:bg-gray-900"
            >
              {submitted ? 'Subscribed!' : 'Subscribe'}
            </button>
          </form>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
