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
  const { loading, error, data } = useQuery(TEST_QUERY);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h2>GraphQL Connection Test</h2>
      <p>Connection successful! Schema types available:</p>
      <ul>
        {data.__schema.types.slice(0, 5).map((type: { name: string }) => (
          <li key={type.name}>{type.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default TestQuery;