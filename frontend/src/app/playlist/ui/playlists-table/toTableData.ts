import { getProviderTextColor } from '@utils/getProviderStyle';

import { Playlist, PlaylistsTableData } from '../../playlist';

const link = ":provider/:playlistId";

export const playlistsToTableData = (
  playlists: Playlist[],
): PlaylistsTableData[] => {
  return playlists.map((playlist, index) => {
    return {
      id: playlist.id,
      link: link
        .replace(":provider", playlist.platform)
        .replace(":playlistId", playlist.id),
      content: [
        {
          type: "select",
          isSelected: playlist.isSelected,
        },
        {
          type: "index",
          index: index,
        },
        {
          type: "title",
          title: playlist.title,
          imageurl: playlist.picture,
        },
        {
          type: "text",
          text: playlist.platform,
          class: getProviderTextColor(playlist.platform),
        },
        {
          type: "text",
          text: `Faixas: ${playlist.total_tracks.toString()}`,
        },
        {
          type: "text",
          text: playlist.creator,
        },
        {
          type: "text",
          text: playlist.public ? "Pública" : "Privada",
        },
      ],
    } as PlaylistsTableData;
  });
};
