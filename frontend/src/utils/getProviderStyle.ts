import { SupportedSources } from "../types/providerAuth";

const providerStyle = {
  spotify: {
    "text-color": "text-[#00F527]",
  },
  youtube: {
    "text-color": "text-[#FF0000]",
  },
};

export const getProviderTextColor = (provider: SupportedSources): string => {
  return providerStyle[provider]["text-color"];
};
