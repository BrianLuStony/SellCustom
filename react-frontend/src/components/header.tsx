import { MainNav } from "./MainNav"
import { SearchBar } from "./SearchBar" // You'll need to create this component
import { UserActions } from "./UserActions" // You'll need to create this component
import { Link } from 'react-router-dom';

export default function Header() {
  return (
    <header className="sticky top-0 w-full bg-white border-b z-50 px-4 py-2">
      <div className="flex items-center justify-between max-w-full mx-52">
        <div className="flex items-center space-x-4">
          <MainNav />
          <SearchBar />
        </div>
        <Link
          to={`/`}
        >
        <div className="flex-grow flex justify-center">
          Logo
        </div>
        </Link>
        <UserActions />
      </div>
    </header>
  )
}