export interface ApiResponse<Data = Record<string, unknown>> {
  count?: number;
  data?: Data;
  error?: string;
}

export async function http<T>(
  url: string,
  {
    method = 'GET',
    body,
    fallbackMessage,
    ...options
  }: RequestInit & { fallbackMessage?: string } = {}
) {
  try {
    const response = await fetch(url, {
      method,
      body,
      ...options,
    });

    const result = (await response.json()) as ApiResponse<NonNullable<T>>;

    if (!response.ok || !result.data) {
      throw new Error(
        result.error || fallbackMessage || 'Failed to perform action.'
      );
    }

    return result;
  } catch (error) {
    console.error(error); // Log the error to an error reporting service
    throw new Error(fallbackMessage || 'Failed to perform action.');
  }
}
