import useFetch, { Options } from 'use-http';

export default function useHttp<T>(endpoint: string, method?: 'POST' | 'PATCH' | 'PUT' | 'DELETE' | 'GET', config?: (options: Options) => Options) {
  const { post, put, patch, del, get, data, loading, error, response } = useFetch<T>(endpoint, options => {
    if (config) {
      config(options);
    }
    return options;
  });

  const fire = <R extends object>(data?: R, headers?: []): void => {
    if (method === 'POST') post(data);
    else if (method === 'PUT') put(data);
    else if (method === 'PATCH') patch(data);
    else if (method === 'DELETE') del();
    else get();
  };

  return {
    fire,
    data,
    loading,
    error: error?.message,
    headers: response.headers
  };
}
