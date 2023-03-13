import { z } from 'zod';
import format from 'date-fns/format';
import { http } from './http';
import { getAPIURL } from '@/utils/env';
import { formatCurrency } from '@/utils';

export const concertSchema = z.object({
  id: z.number(),
  title: z.string(),
  subtitle: z.string(),
  date: z.coerce.date().transform(date => ({
    value: date,
    formatted: format(date, 'LLLL dd, yyyy'),
    formattedStartTime: format(date, 'HH:mmaaa'),
  })),
  ticket_price: z.number().transform(price => ({
    value: price,
    formatted: formatCurrency({ amount: price / 100 }),
  })),
  venue: z.string(),
  venue_address: z.string(),
  city: z.string(),
  state: z.string(),
  zip: z.string(),
  additional_information: z.string(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
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
  const { data } = await http<unknown[]>(`${getAPIURL()}/concerts`, {
    fallbackMessage: 'Could not retrieve concerts',
  });

  return parseConcerts(data);
}
