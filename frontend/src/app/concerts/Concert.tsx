import {
  CalendarDaysIcon,
  ClockIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/24/solid';

function Concert() {
  return (
    <article className="md :w-6/12 w-full space-y-8 rounded-xl bg-white p-10  shadow-md">
      <div className="space-y-2">
        <h3 className="text-3xl font-bold capitalize">The Red Chord</h3>
        <p className="font-semibold">With Animosity and Lethargy</p>
      </div>

      <div className="space-y-6">
        <div className="flex items-start gap-4">
          <CalendarDaysIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">December 13, 2016</p>
        </div>
        <div className="flex items-start gap-4">
          <ClockIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">Doors at 8:00pm</p>
        </div>
        <div className="flex items-start gap-4">
          <CurrencyDollarIcon className="h-6 w-6 text-gray-600" />
          <p className="font-semibold text-gray-700">32.50</p>
        </div>
        <div className="flex items-start gap-4">
          <MapPinIcon className="h-6 w-6 text-gray-600" />
          <div className="space-y-2">
            <p className="font-semibold text-gray-700">The Mosh Pit</p>
            <p>123 Example Lane</p>
            <p>Laravile, ON 17916</p>
          </div>
        </div>
        <div className="flex items-start gap-4">
          <ExclamationCircleIcon className="h-6 w-6 text-gray-600" />
          <div className="space-y-2">
            <p className="font-semibold text-gray-700">Addtional Information</p>
            <p>For tickets, call (555) 555-5555</p>
          </div>
        </div>
      </div>
    </article>
  );
}

export default Concert;
