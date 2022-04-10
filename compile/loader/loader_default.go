package loader

/*

// Loader which loads defines.RawPost
type RawPost struct {
	Loader fsal.DataLoader
}

func (bc *RawPost) Load(ctx context.Context, fs LoaderInput) (iterator iter.Iterable[defines.RawPost], err error) {
	iterator = iter.IterableFunc[defines.RawPost](func(ctx context.Context, recv iter.Receiver[defines.RawPost]) (err error) {
		entries, err := fs.ReadDir("/")
		if err != nil {
			return
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			postDir := e.Name()
			if strings.HasPrefix(postDir, ".") || strings.HasPrefix(postDir, "_") {
				continue
			}

			var meta defines.RawPostMetadata
			var content defines.PostContent

			err = bc.Loader.ReadData(fs, path.Join(postDir, "metadata"), &meta)
			if err != nil {
				return
			}

			err = bc.Loader.ReadData(fs, path.Join(postDir, "content"), &content)
			if err != nil {
				return
			}

			var post defines.RawPost
			post.Metadata = meta
			post.Dir = postDir
			post.Content = content

			post.FS = &fsal.PrefixFS{
				Wrapped:    fs,
				PathPrefix: postDir,
			}

			err = recv(ctx, post)
			if err != nil {
				return
			}
		}

		return
	})

	return
}
*/
