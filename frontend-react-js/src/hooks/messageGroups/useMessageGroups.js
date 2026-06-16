import { useQuery } from "@tanstack/react-query";
import { fetchMessageGroup } from "../../api/messageGroups";

export function useMessageGroup(handle, token) {
  return useQuery({
    queryKey: ["message-group", handle, token],
    queryFn: () => fetchMessageGroup(handle, token),
    enabled: !!handle && !!token,
  });
}