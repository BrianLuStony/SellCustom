import React, { useState } from 'react';
import { gql, useMutation } from '@apollo/client';

const CREATE_PRODUCT = gql`
  mutation CreateProduct($input: ProductInput!) {
    createProduct(input: $input) {
      id
      name
      description
      price
      category {
        id
        name
      }
    }
  }
`;

const AddProduct: React.FC = () => {
  const [createProduct, { data, loading, error }] = useMutation(CREATE_PRODUCT);
  const [productInput, setProductInput] = useState({
    name: '',
    description: '',
    price: '',
    stockQuantity: '',
    categoryId: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProductInput(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const input = {
      ...productInput,
      price: parseFloat(productInput.price),
      stockQuantity: parseInt(productInput.stockQuantity, 10),
    };
    createProduct({ 
      variables: { 
        input: {
          ...input,
          stockQuantity: input.stockQuantity || 0, // Ensure it's a number, default to 0 if NaN
        } 
      } 
    });
  };

  if (loading) return <p>Submitting...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h1>Add New Product</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="name"
          placeholder="Product Name"
          value={productInput.name}
          onChange={handleChange}
        />
        <input
          type="text"
          name="description"
          placeholder="Description"
          value={productInput.description}
          onChange={handleChange}
        />
        <input
          type="number"
          name="price"
          placeholder="Price"
          value={productInput.price}
          onChange={handleChange}
          step="0.01"
        />
        <input
          type="number"
          name="stockQuantity"
          placeholder="Stock Quantity"
          value={productInput.stockQuantity}
          onChange={handleChange}
          step="1"
        />
        <input
          type="text"
          name="categoryId"
          placeholder="Category ID"
          value={productInput.categoryId}
          onChange={handleChange}
        />
        <button type="submit">Add Product</button>
      </form>
      {data && <p>Product added: {data.createProduct.name}</p>}
    </div>
  );
};

export default AddProduct;