import { getAuthUser } from "@/auth/server/action";
import { Suspense } from "react";
import { SearchContent } from "./search-content";

export default async function SearchPage() {
  const auth = await getAuthUser();
  return (
    <Suspense
      fallback={
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-muted rounded w-1/3 mx-auto"></div>
            <div className="h-32 bg-muted rounded"></div>
          </div>
        </div>
      }
    >
      <SearchContent auth={auth} />
    </Suspense>
  );
}
