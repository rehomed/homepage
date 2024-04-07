import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "./index.css";
import Main from "./routes/Main";
import { LocationProvider, ErrorBoundary, Route, Router } from "preact-iso";

const queryClient = new QueryClient();

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ErrorBoundary>
        <LocationProvider>
          <Router>
            <Route component={Main} default />
            <Route component={Main} default />
          </Router>
        </LocationProvider>
      </ErrorBoundary>
    </QueryClientProvider>
  );
}
