import React from 'react';
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

const Products: React.FC = () => {
  const { loading, error, data } = useQuery<{ products: Product[] }>(GET_PRODUCTS);

  if (loading) return <p>Loading products...</p>;
  if (error) return <p>Error loading products: {error.message}</p>;
  if (!data || !data.products) return <p>No product data available.</p>;

  return (
    <div>
      <h1>Products</h1>
      {data.products.length > 0 ? (
        <div className="product-grid">
          {data.products.map((product) => (
            <div key={product.id} className="product-card">
              <h2>{product.name}</h2>
              <p>{product.description}</p>
              <p>Price: ${product.price.toFixed(2)}</p>
              <p>Category: {product.category.name}</p>
              <div className="product-images">
                {product.images.map((image) => (
                  <img key={image.id} src={image.imageUrl} alt={product.name} />
                ))}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p>No products found.</p>
      )}
    </div>
  );
};

export default Products;
