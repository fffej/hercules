/*
Package hercules contains the functions which are needed to gather various statistics
from a Git repository.

The analysis is expressed in a form of the tree: there are nodes - "pipeline items" - which
require some other nodes to be executed prior to selves and in turn provide the data for
dependent nodes. There are several service items which do not produce any useful
statistics but rather provide the requirements for other items. The top-level items
are:

- BurndownAnalysis - line burndown statistics for project, files and developers.
- Couples - coupling statistics for files and developers.

The typical API usage is to initialize the Pipeline class:

  import "gopkg.in/src-d/go-git.v4"

	var repository *git.Repository
	// ...initialize repository...
	pipeline := hercules.NewPipeline(repository)

Then add the required analysis tree nodes:

  pipeline.AddItem(&hercules.BlobCache{})
	pipeline.AddItem(&hercules.DaysSinceStart{})
	pipeline.AddItem(&hercules.TreeDiff{})
	pipeline.AddItem(&hercules.FileDiff{})
	pipeline.AddItem(&hercules.RenameAnalysis{SimilarityThreshold: 80})
	pipeline.AddItem(&hercules.IdentityDetector{})

Then initialize BurndownAnalysis:

  burndowner := &hercules.BurndownAnalysis{
    Granularity:  30,
		Sampling:     30,
  }
  pipeline.AddItem(burndowner)

Then execute the analysis tree:

  pipeline.Initialize()
	result, err := pipeline.Run(commits)

Finally extract the result:

  burndownResults := result[burndowner].(hercules.BurndownResult)

The actual usage example is cmd/hercules/main.go - the command line tool's code.

Hercules depends heavily on https://github.com/src-d/go-git and leverages the
diff algorithm through https://github.com/sergi/go-diff.

Besides, hercules defines File and RBTree. These are low level data structures
required by BurndownAnalysis. File carries an instance of RBTree and the current line
burndown state. RBTree implements the red-black balanced binary tree and is
based on https://github.com/yasushi-saito/rbtree.

Coupling stats are supposed to be further processed rather than observed directly.
labours.py uses Swivel embeddings and visualises them in Tensorflow Projector.
*/
package hercules
