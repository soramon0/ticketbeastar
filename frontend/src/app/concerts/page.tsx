import { getConcerts } from '@/api/concerts';
import Concert from './Concert';

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
