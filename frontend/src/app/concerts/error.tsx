'use client'; // Error components must be Client components

export default function Error({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  return (
    <main className="container mx-auto flex h-screen items-center justify-center p-4">
      <div className="w-full space-y-8 rounded-xl bg-white p-10 shadow-md md:w-6/12">
        <h1 className="text-3xl font-bold">Something went wrong!</h1>
        <p className="text-xl font-semibold text-red-400">{error.message}</p>
        <button
          className="inline-flex justify-center rounded-md bg-indigo-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
          onClick={() => reset()}
        >
          Try again
        </button>
      </div>
    </main>
  );
}
