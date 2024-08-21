import React, { useState } from 'react';
import { Link } from 'react-router-dom';

interface SidebarProps {
    isOpen: boolean;
    onClose: () => void;
}

const menuItems = [
  { title: 'New In', href: '/new-in' },
  { title: 'On Sale', href: '/sale' },
  {
    title: 'Women',
    href: '/women',
    subItems: [
      { title: 'Tops', href: '/women/tops' },
      { title: 'Bottoms', href: '/women/bottoms' },
      { title: 'Dresses', href: '/women/dresses' },
      { title: 'Accessories', href: '/women/accessories' },
    ],
  },
  {
    title: 'Men',
    href: '/men',
    subItems: [
      { title: 'Tops', href: '/men/tops' },
      { title: 'Bottoms', href: '/men/bottoms' },
      { title: 'Shoes', href: '/men/shoes' },
      { title: 'Accessories', href: '/men/accessories' },
    ],
  },
  { title: 'About', href: '/about' },
];

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const [expandedMenu, setExpandedMenu] = useState<string | null>(null);

  const toggleExpand = (title: string) => {
    setExpandedMenu(expandedMenu === title ? null : title);
  };

  return (
    <>
      {isOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-40" onClick={onClose}></div>
      )}
      <div
        className={`fixed top-0 left-0 h-full w-96 bg-white transform ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } transition-transform duration-300 ease-in-out z-50 flex flex-col justify-between`}
      >
        <div>
          <div className="flex justify-between items-center p-4 border-b">
            <h2 className="text-xl font-semibold">Menu</h2>
            <button onClick={onClose} className="text-2xl">&times;</button>
          </div>
          <nav className="p-4 flex-1 overflow-y-auto">
            <ul className="space-y-2">
              {menuItems.map((item) => (
                <li key={item.title}>
                  <div
                    className="flex justify-between items-center cursor-pointer py-2 px-4 text-gray-700 hover:bg-gray-100 rounded"
                    onClick={() => toggleExpand(item.title)}
                  >
                    <Link
                      to={`${item.href.toLowerCase().replace(/ /g, '-')}`}
                      className="block py-2 px-4 text-gray-600 rounded"
                      onClick={onClose}
                    >
                      {item.title}
                    </Link>
                    {item.subItems && (
                      <span className="text-gray-500">
                        {expandedMenu === item.title ? '-' : '+'}
                      </span>
                    )}
                  </div>
                  {expandedMenu === item.title && item.subItems && (
                    <ul className="ml-6 space-y-1">
                      {item.subItems.map((subItem, index) => (
                        <li key={index}>
                          <Link
                            to={`${item.href}/${subItem.href.toLowerCase().replace(/ /g, '-')}`}
                            className="block py-2 px-4 text-gray-600 hover:bg-gray-200 rounded"
                            onClick={onClose}
                          >
                            {subItem.title}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  )}
                </li>
              ))}
            </ul>
          </nav>
        </div>
        <div className="p-4 border-t space-y-2 bg-zinc-400">
          <Link
            to="/account"
            className="block py-2 pl-4 text-gray-600 hover:bg-gray-300 rounded text-left"
            onClick={onClose}
          >
            Account
          </Link>
          <Link
            to="/contact"
            className="block py-2 pl-4 text-gray-600 hover:bg-gray-300 rounded text-left"
            onClick={onClose}
          >
            Contact
          </Link>
        </div>
      </div>
    </>
  );
};

export default Sidebar;
