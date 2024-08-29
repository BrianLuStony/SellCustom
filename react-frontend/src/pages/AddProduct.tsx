import React, { useState } from 'react';
import { gql, useMutation } from '@apollo/client';
import * as UploadcareWidget from '@uploadcare/react-widget';

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
  const [createProduct, { loading, error }] = useMutation(CREATE_PRODUCT, {
    onError: (error) => {
      console.error('Error creating product:', error);
    },
    onCompleted: () => {
      setSuccessMessage('Product added successfully!');
      resetForm();
    },
  });

  const [productInput, setProductInput] = useState({
    name: '',
    description: '',
    price: '',
    stockQuantity: '',
    categoryId: '',
    images: [] as { imageUrl: string; isPrimary: boolean }[], // Initialize with correct type
  });

  const [successMessage, setSuccessMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProductInput(prev => ({ ...prev, [name]: value }));
  };

  const handleImageUpload = (fileInfo :any) => {
    setProductInput(prev => ({
      ...prev,
      images: [...prev.images, { imageUrl: fileInfo.cdnUrl, isPrimary: true }],
    }));
  };

  const resetForm = () => {
    setProductInput({
      name: '',
      description: '',
      price: '',
      stockQuantity: '',
      categoryId: '',
      images: [],
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage('');
    const input = {
      ...productInput,
      price: parseFloat(productInput.price),
      stockQuantity: parseInt(productInput.stockQuantity, 10) || 0,
      categoryId: parseInt(productInput.categoryId, 10),
    };
    createProduct({ variables: { input } });
  };

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

        {/* Uploadcare Widget for Image Upload */}
        <UploadcareWidget.Widget
          publicKey="3cb90b7ae33a12d4d266" // Replace with your actual public key
          onChange={handleImageUpload}
        />

        <button type="submit" disabled={loading}>
          {loading ? 'Adding...' : 'Add Product'}
        </button>
      </form>
      {loading && <p>Submitting...</p>}
      {error && <p>Error: {error.message}</p>}
      {successMessage && <p>{successMessage}</p>}
    </div>
  );
};

export default AddProduct;
