import React, { useState } from 'react';
import { gql, useMutation } from '@apollo/client';

const CREATE_CATEGORY = gql`
  mutation CreateCategory($input: CategoryInput!) {
    createCategory(input: $input) {
      id
      name
      parentCategory {
        id
        name
      }
    }
  }
`;

const AddCategory: React.FC = () => {
  const [createCategory, { data, loading, error }] = useMutation(CREATE_CATEGORY);
  const [categoryInput, setCategoryInput] = useState({
    name: '',
    parentId: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCategoryInput(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const input = {
      ...categoryInput,
      parentId: categoryInput.parentId ? categoryInput.parentId : null,
    };
    createCategory({ variables: { input } });
  };

  if (loading) return <p>Submitting...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h1>Add New Category</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="name"
          placeholder="Category Name"
          value={categoryInput.name}
          onChange={handleChange}
          required
        />
        <input
          type="text"
          name="parentId"
          placeholder="Parent Category ID (optional)"
          value={categoryInput.parentId}
          onChange={handleChange}
        />
        <button type="submit">Add Category</button>
      </form>
      {data && <p>Category added: {data.createCategory.name}</p>}
    </div>
  );
};

export default AddCategory;