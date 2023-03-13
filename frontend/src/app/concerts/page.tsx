import Concert from './Concert';
import { parseConcerts } from '@/api/concerts';

interface ApiResponse<Data = Record<string, unknown>> {
  count?: number;
  data?: Data;
  error?: string;
}

async function getConcerts() {
  const response = await fetch(`http://localhost:5000/api/v1/concerts`);
  const result = (await response.json()) as ApiResponse<unknown[]>;

  if (!response.ok || !result.data) {
    throw new Error(result.error || 'Could not retrieve concerts');
  }

  return parseConcerts(result.data);
}

async function Home() {
  const concerts = await getConcerts();

  return (
    <main className="container mx-auto flex h-screen items-center justify-center p-4">
      {concerts.map(concert => (
        <Concert key={concert.id} concert={concert} />
      ))}
    </main>
  );
}

export default Home;
