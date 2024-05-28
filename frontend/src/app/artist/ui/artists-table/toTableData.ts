import { Artist, ArtistsTableData } from 'app/artist/artist';

import { formatNumber } from '@utils/formatNumber';
import { getProviderTextColor } from '@utils/getProviderStyle';

export const artistsToTableData = (artists: Artist[]): ArtistsTableData[] => {
  return artists.map((album, index) => {
    return {
      id: album.id,
      content: [
        {
          type: "select",
          isSelected: album.isSelected,
        },
        {
          type: "index",
          index: index,
        },
        {
          type: "title",
          title: album.name,
          imageurl: album.picture,
          link: album.link,
        },
        {
          type: "text",
          text: album.platform,
          class: getProviderTextColor(album.platform),
          link: album.link,
        },
        {
          type: "text",
          text: `Fãs: ${formatNumber(album.fans)}`,
          link: album.link,
        },
      ],
    } as ArtistsTableData;
  });
};
