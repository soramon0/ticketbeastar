import {
  CalendarDaysIcon,
  ClockIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/24/solid';
import type { Concert as IConcert } from '@/api/concerts';
import { formatCurrency } from '@/utils';
import { format } from 'date-fns';

interface Props {
  concert: IConcert;
}

function Concert({ concert }: Props) {
  return (
    <article className="md :w-6/12 w-full space-y-8 rounded-xl bg-white p-10  shadow-md">
      <div className="space-y-2">
        <h3 className="text-3xl font-bold capitalize">{concert.title}</h3>
        <p className="font-semibold">{concert.subtitle}</p>
      </div>

      <div className="space-y-6">
        <div className="flex items-start gap-4">
          <CalendarDaysIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">
            {format(concert.date, 'LLLL dd, yyyy')}
          </p>
        </div>
        <div className="flex items-start gap-4">
          <ClockIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">
            Doors at {format(concert.date, 'HH:mmaaa')}
          </p>
        </div>
        <div className="flex items-start gap-4">
          <CurrencyDollarIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">
            {formatCurrency({ amount: concert.ticket_price / 100 })}
          </p>
        </div>
        <div className="flex items-start gap-4">
          <MapPinIcon className="h-6 w-6 text-gray-600" />
          <div className="space-y-2">
            <p className="font-semibold text-gray-700">{concert.venue}</p>
            <p>{concert.venue_address}</p>
            <p>
              {concert.city}, {concert.state} {concert.zip}
            </p>
          </div>
        </div>
        <div className="flex items-start gap-4">
          <ExclamationCircleIcon className="h-6 w-6 text-gray-600" />
          <div className="space-y-2">
            <p className="font-semibold text-gray-700">Addtional Information</p>
            <p>{concert.additional_information}</p>
          </div>
        </div>
      </div>
    </article>
  );
}

export default Concert;
