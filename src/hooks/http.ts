import { useState } from 'react';

export default function useHttp<T>(endpoint: string, method?: 'POST' | 'PATCH' | 'PUT' | 'DELETE' | 'GET') {
  const [data, setData] = useState<T>();
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<any>();
  const [responseHeaders, setResponseHeaders] = useState<Headers>();

  const fire = (body?: any, json: boolean = true, headers: Headers | string[][] | Record<string, string> = {}): void => {
    setLoading(true);
    const request: RequestInit = {
      method,
      body: json ? JSON.stringify(body) : body,
      credentials: 'include',
      headers: json ? { 'Content-Type': 'application/json', ...headers } : headers
    };
    fetch(endpoint, request).then(res => {
      setResponseHeaders(res.headers);
      return res.json();
    }).then(json => {
      setData(json);
      setLoading(false);
    }).catch(err => {
      setData(undefined);
      setError(err);
      setLoading(false);
    });
  };


  return {
    fire,
    data,
    loading,
    error,
    headers: responseHeaders
  };
}
