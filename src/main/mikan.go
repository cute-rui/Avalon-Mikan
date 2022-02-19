package main

import (
	"avalon-mikan/src/mikan"
	"context"
	"github.com/mmcdole/gofeed"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type MikanServiceServer struct {
	mikan.UnimplementedMikanServer
}

const (
	MIKAN_TORRENT = 0
	MIKAN_MAGLINK = 1
)

func (m MikanServiceServer) GetInfo(ctx context.Context, param *mikan.Param) (*mikan.Info, error) {
	BangumiID, SubgroupID := param.GetBangumi(), param.GetSubgroup()

	if BangumiID == `` {
		return nil, status.Error(codes.InvalidArgument, `BangumiID is empty`)
	}

	if SubgroupID == `` {
		return nil, status.Error(codes.Unimplemented, `Unimplemented method`)
	}

	query := StringBuilder(`bangumiId=`, BangumiID)
	if SubgroupID != `` {
		query = StringBuilder(query, `&subgroupid=`, SubgroupID)
	}

	URL := StringBuilder(`https://mikanani.me/RSS/Bangumi?`, query)

	feedParser := gofeed.NewParser()

	feed, err := feedParser.ParseURL(URL)
	if err != nil {
		return nil, status.Error(400, err.Error())
	}

	info := mikan.Info{
		BangumiName: feed.Title,
		BangumiID:   BangumiID,
		SubgroupID:  SubgroupID,
	}

	for i := range feed.Items {
		var d mikan.Data
		d.Title = feed.Items[i].Title
		d.URL = feed.Items[i].Enclosures[0].URL
		d.DataType = MIKAN_TORRENT

		info.Data = append(info.GetData(), &d)
	}

	return &info, nil
}

func StringBuilder(p ...string) string {
	var (
		b strings.Builder
		c int
	)
	l := len(p)
	for i := 0; i < l; i++ {
		c += len(p[i])
	}
	b.Grow(c)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}
