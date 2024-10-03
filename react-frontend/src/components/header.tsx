import { MainNav } from "./MainNav"
import { SearchBar } from "./SearchBar" // You'll need to create this component
import { UserActions } from "./UserActions" // You'll need to create this component
import { Link } from 'react-router-dom';

interface HeaderProps {
  onSearch: (term: string) => void;
}

export default function Header({ onSearch }: HeaderProps) {
  return (
    <header className="sticky top-0 w-full bg-white border-b z-50 px-4 py-2">
      <div className="flex items-center justify-between max-w-full mx-52">
        <div className="flex items-center space-x-4">
          <MainNav />
          {/* Pass the onSearch function to SearchBar */}
          <SearchBar onSearch={onSearch} />
        </div>
        <Link to={`/`}>
          <div className="flex-grow flex justify-center">Logo</div>
        </Link>
        <UserActions />
      </div>
    </header>
  );
}