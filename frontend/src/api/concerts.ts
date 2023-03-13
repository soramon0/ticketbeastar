import { z } from 'zod';
import { ApiResponse } from './http';

export const concertSchema = z.object({
  id: z.number(),
  title: z.string(),
  subtitle: z.string(),
  date: z.string().transform(date => new Date(date)),
  ticket_price: z.number(),
  venue: z.string(),
  venue_address: z.string(),
  city: z.string(),
  state: z.string(),
  zip: z.string(),
  additional_information: z.string(),
  created_at: z.string().transform(date => new Date(date)),
  updated_at: z.string().transform(date => new Date(date)),
});

export type IConcert = z.infer<typeof concertSchema>;

export function parseConcert(data: unknown): IConcert {
  return concertSchema.parse(data);
}

export function parseConcerts(data: unknown): IConcert[] {
  if (!Array.isArray(data)) {
    throw new Error('concert data is not a list.');
  }

  return data.map(parseConcert);
}

export async function getConcerts() {
  const response = await fetch(`http://localhost:5000/api/v1/concerts`);
  const result = (await response.json()) as ApiResponse<unknown[]>;

  if (!response.ok || !result.data) {
    throw new Error(result.error || 'Could not retrieve concerts');
  }

  return parseConcerts(result.data);
}
