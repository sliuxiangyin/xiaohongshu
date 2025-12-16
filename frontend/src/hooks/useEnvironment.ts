import { useEffect, useState } from 'react';
import { Environment } from '../../wailsjs/runtime';

export interface EnvironmentInfo {
  buildType: string;
  platform: string;
  arch: string;
}

export const useEnvironment = () => {
  const [envInfo, setEnvInfo] = useState<EnvironmentInfo | null>(null);
  const [isDev, setIsDev] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchEnvironment = async () => {
      try {
        setLoading(true);
        const environment = await Environment();
        setEnvInfo(environment);
        // 在 Wails 中，buildType 为 'dev' 表示开发环境，其他值表示生产环境
        setIsDev(environment.buildType === 'dev');
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to get environment info'));
        console.error('Failed to get environment info:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchEnvironment();
  }, []);

  return { envInfo, isDev, loading, error };
};