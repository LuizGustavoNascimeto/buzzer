import { useQuery } from "@tanstack/react-query";
import { fetchMessages } from "../../api/messages";
import { useAuth } from "../auth/useAuth";

export function useMessage(group_id, token) {
  return useQuery({
    queryKey: ["message", group_id],
    queryFn: () => fetchMessages(group_id, token),
    enabled: !!group_id && !!token,
  });
}