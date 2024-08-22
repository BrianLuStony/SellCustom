import React from 'react';
import { useQuery, gql } from '@apollo/client';

const TEST_QUERY = gql`
  query {
    __schema {
      types {
        name
      }
    }
  }
`;

const TestQuery: React.FC = () => {
    const { loading, error, data } = useQuery(TEST_QUERY, {
        fetchPolicy: 'network-only', // This ensures we always make a network request
      });

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;
      console.log(data);
  return (
    <div>
      <h2>GraphQL Connection Test</h2>
      <p>Connection successful! Schema types available:</p>
      <ul>
        {data.__schema.types.map((type: { name: string }) => (
          <li key={type.name}>{type.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default TestQuery;