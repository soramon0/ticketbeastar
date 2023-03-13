import { z } from 'zod';

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

export type Concert = z.infer<typeof concertSchema>;

export function parseConcert(data: unknown): Concert {
  return concertSchema.parse(data);
}

export function parseConcerts(data: unknown): Concert[] {
  if (!Array.isArray(data)) {
    throw new Error('concert data is not a list.');
  }

  return data.map(parseConcert);
}
