// components/Sidebar.jsx
import React from 'react';
import { Link } from 'react-router-dom';

interface SidebarProps {
    isOpen: boolean;
    onClose: () => void;
}
const menuItems = [
  { title: 'Women', href: '/women' },
  { title: 'Men', href: '/men' },
  { title: 'Home, outdoor and equestrian', href: '/home-outdoor-equestrian' },
  { title: 'Jewelry and watches', href: '/jewelry-watches' },
  { title: 'Fragrances and make-up', href: '/fragrances-makeup' },
  { title: 'Gifts and Petit H', href: '/gifts-petit-h' },
  { title: 'Special editions and services', href: '/special-editions-services' },
  { title: 'About Hermès', href: '/about' },
];

const Sidebar:React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  return (
    <>
      {isOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-40" onClick={onClose}></div>
      )}
      <div
        className={`fixed top-0 left-0 h-full w-64 bg-white transform ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } transition-transform duration-300 ease-in-out z-50`}
      >
        <div className="flex justify-between items-center p-4 border-b">
          <h2 className="text-xl font-semibold">Menu</h2>
          <button onClick={onClose} className="text-2xl">&times;</button>
        </div>
        <nav className="p-4">
          <ul className="space-y-2">
            {menuItems.map((item) => (
              <li key={item.title}>
                <Link
                  to={item.href}
                  className="block py-2 px-4 text-gray-700 hover:bg-gray-100 rounded"
                  onClick={onClose}
                >
                  {item.title}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>
    </>
  );
};

export default Sidebar;