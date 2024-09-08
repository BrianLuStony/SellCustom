import Products from "./Products";

interface HomeProps {
  searchTerm: string;
}

function Home({ searchTerm }: HomeProps) {
  return (
    <div className="min-h-screen bg-gray-100">
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold text-gray-900">Home Page</h1>
        {/* Pass the searchTerm to the Products component */}
        <Products searchTerm={searchTerm} />
      </main>
    </div>
  );
}

export default Home;
