import { useState } from "react";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export default function useHttp<T>(endpoint: string, method?: "POST" | "PATCH" | "PUT" | "DELETE" | "GET") {
  const [data, setData] = useState<T>();
  const [loading, setLoading] = useState(false);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [error, setError] = useState<any>();
  const [status, setStatus] = useState(0);
  const [ok, setOk] = useState(true);
  const [responseHeaders, setResponseHeaders] = useState<Headers>();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const fire = (body?: any, json = true, overwriteUrl?: string, headers: HeadersInit = {}): void => {
    setLoading(true);
    setOk(false);
    setStatus(0);
    const request: RequestInit = {
      method,
      body: json ? JSON.stringify(body) : body,
      credentials: "include",
      headers: json ? { "Content-Type": "application/json", ...headers } : headers
    };
    fetch(overwriteUrl ? overwriteUrl : endpoint, request).then(res => {
      setResponseHeaders(res.headers);
      setOk(res.ok);
      setStatus(res.status);
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
    headers: responseHeaders,
    status,
    ok
  };
}
