import { useQuery, gql } from '@apollo/client';

const GET_PRODUCTS = gql`
  query GetProducts {
    products {
      id
      name
      description
      price
      category {
        id
        name
      }
      images {
        id
        imageUrl
        isPrimary
      }
    }
  }
`;

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  category: {
    id: string;
    name: string;
  };
  images: {
    id: string;
    imageUrl: string;
    isPrimary: boolean;
  }[];
}

interface ProductsProps {
  searchTerm: string;
}

const Products: React.FC<ProductsProps> = ({ searchTerm }) => {
  const { loading, error, data } = useQuery(GET_PRODUCTS);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error loading products</p>;

  const filteredProducts = data.products.filter((product: Product) =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.category.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div>
      {filteredProducts.length === 0 ? (
        <p>No products found</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {filteredProducts.map((product: Product) => (
            <div key={product.id} className="p-4 border rounded-lg shadow-md">
              <h2 className="text-xl font-bold">{product.name}</h2>
              <p>{product.description}</p>
              <p>Price: ${product.price.toFixed(2)}</p>
              {product.images.map((image) => (
                <img
                  key={image.id}
                  src={image.imageUrl}
                  alt={product.name}
                  className="w-full h-auto"
                />
              ))}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Products;
