import Products from "./Products";
function Home() {
  return (
    <div className="min-h-screen bg-gray-100">
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold text-gray-900">Home Page</h1>
        <Products />
        {/* Add your main content here */}
      </main>
    </div>
  );
}

export default Home;