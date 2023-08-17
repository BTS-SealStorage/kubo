package commands

import (
	"fmt"
	cmdenv "github.com/ipfs/kubo/core/commands/cmdenv"
	"io"
	"net/url"

	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/files"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var s3StoreCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with s3store.",
	},
	Subcommands: map[string]*cmds.Command{
		"add": s3Add,
	},
}

var s3Add = &cmds.Command{
	//Status: cmds.Deprecated,
	Helptext: cmds.HelpText{
		Tagline: "Add S3 object via urlstore.",
		LongDescription: `
Add S3 objects to ipfs without storing the data locally.

The S3 URL provided must be stable and ideally under your
control.

The file is added using raw-leaves but otherwise using the default
settings for 'ipfs add'.
`,
	},
	Options: []cmds.Option{
		cmds.BoolOption(trickleOptionName, "t", "Use trickle-dag format for dag generation."),
		cmds.BoolOption(pinOptionName, "Pin this object when adding.").WithDefault(true),
	},
	Arguments: []cmds.Argument{
		//cmds.StringArg("s3endpoint", true, false, "S3 endpoint URL to find object"),
		cmds.StringArg("url", true, false, "object URL to add to IPFS"),
	},
	Type: &BlockStat{},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		//log.Error("The 'ipfs urlstore' command is deprecated, please use 'ipfs add --nocopy --cid-version=1")

		urlString := req.Arguments[0]
		if !IsS3URL(urlString) {
			return fmt.Errorf("unsupported url syntax: %s", urlString)
		}

		url, err := url.Parse(urlString)
		if err != nil {
			return err
		}

		enc, err := cmdenv.GetCidEncoder(req)
		if err != nil {
			return err
		}

		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		useTrickledag, _ := req.Options[trickleOptionName].(bool)
		dopin, _ := req.Options[pinOptionName].(bool)

		opts := []options.UnixfsAddOption{
			options.Unixfs.Pin(dopin),
			options.Unixfs.CidVersion(1),
			options.Unixfs.RawLeaves(true),
			options.Unixfs.Nocopy(true),
		}

		if useTrickledag {
			opts = append(opts, options.Unixfs.Layout(options.TrickleLayout))
		}

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		file := files.NewS3File(&node.S3Connection, url)

		path, err := api.Unixfs().Add(req.Context, file, opts...)
		if err != nil {
			return err
		}
		size, _ := file.Size()
		return cmds.EmitOnce(res, &BlockStat{
			Key:  enc.Encode(path.Cid()),
			Size: int(size),
		})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, bs *BlockStat) error {
			_, err := fmt.Fprintln(w, bs.Key)
			return err
		}),
	},
}

// IsS3URL returns true if the string represents a valid S3 endpoint
// that the urlstore can handle.  More specifically it returns true
// if a string begins with 's3://'.
func IsS3URL(str string) bool {
	return len(str) > 5 && str[0] == 's' && str[1] == '3' && str[2] == ':' && str[3] == '/' && str[4] == '/'
}
