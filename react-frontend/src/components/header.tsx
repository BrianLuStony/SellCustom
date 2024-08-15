import { MainNav } from "./main-nav"
// import UserButton from "./Sign-components/user-button"
export default function Header() {
  return (
    <header className="sticky top-0 flex justify-center border-b bg-slate-800 z-50 dark:bg-slate-800 px-4 py-2">
      <div className="flex items-center justify-between w-full h-16 max-w-7xl mx-auto">
        <MainNav />
        {/* <UserButton /> */}
      </div>
    </header>
  )
}
