import { useState, useEffect, useRef } from "react";
import { ApiStatus } from "../types";

type HealthCheckFunction = () => Promise<any>;

export const useApiStatus = (
  healthCheckFn: HealthCheckFunction,
  apiName: string
): ApiStatus | null => {
  const [apiStatus, setApiStatus] = useState<ApiStatus | null>(null);
  const healthCheckFnRef = useRef(healthCheckFn);

  // Update ref when function changes
  useEffect(() => {
    healthCheckFnRef.current = healthCheckFn;
  }, [healthCheckFn]);

  useEffect(() => {
    const checkApiStatus = async (): Promise<void> => {
      try {
        const healthResponse = await healthCheckFnRef.current();
        
        // Handle different response formats
        let message = "";
        if (typeof healthResponse === "string") {
          message = healthResponse;
        } else if (healthResponse.message) {
          message = healthResponse.message;
        } else if (healthResponse.status) {
          message = `${apiName} is running`;
        } else {
          message = `${apiName} is online`;
        }
        
        setApiStatus({
          status: "online",
          message: message,
        });
      } catch (error: any) {
        setApiStatus({
          status: "error",
          message: `Failed to connect to ${apiName}`,
        });
      }
    };

    checkApiStatus();
  }, [apiName]);

  return apiStatus;
};

