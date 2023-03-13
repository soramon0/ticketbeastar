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
  const response = await fetch(url, {
    method,
    body,
    ...options,
  });

  const { data, error } = (await response.json()) as ApiResponse<T>;

  if (!response.ok || !data) {
    throw new Error(error || fallbackMessage || 'Failed to perform action.');
  }

  return data;
}
