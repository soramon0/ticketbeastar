function checkEnv(env: string | undefined, name: string) {
  if (!env) {
    throw new Error(
      `Please define the ${name} environment variable inside .env`
    );
  }

  return env;
}

export function getAPIURL() {
  const value = process.env.NEXT_PUBLIC_API_URL;
  return checkEnv(value, 'NEXT_PUBLIC_API_URL');
}
