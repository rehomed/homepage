import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "./index.css";
import Main from "./main";

const queryClient = new QueryClient();

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Main />
    </QueryClientProvider>
  );
}
