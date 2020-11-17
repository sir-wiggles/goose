package lib

import (
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to the last marker",
	Long: `Rollbacks will rollback a batch of migrations using the marker talked about above. 
For exampe, here we have two batches:  
	* 1: a, b and c 
	* 2: d, e and f
Where f and c are markers indicating the last migration ran in their batch.

+------------------------------------------+----------+----------+
| hash                                     | author   | marker   |
|------------------------------------------+----------+----------|
| e965f4511fce6ae61e1cfdcf174f61cfd4fe920b | a o      | False    |
| cac4966fa648df678b9f59117d085b40d647ef19 | b o      | False    |
| e0ca0a9d0afe2d168ed09efe2f859f76bcfd109f | c o      | True     |
| 6a8f40ecd57b264da0d0492af62b577f626bfbe1 | d o      | False    |
| 76499a490b9c0006100d963e6006f72cf56c6826 | e o      | False    |
| 9ebb39681a4428cc5693ea2d926e5f73711ce9a4 | f o      | True     |
+------------------------------------------+----------+----------+

To rollback to c run "goose rollback" which will put us in this state

+------------------------------------------+----------+----------+
| hash                                     | author   | marker   |
|------------------------------------------+----------+----------|
| e965f4511fce6ae61e1cfdcf174f61cfd4fe920b | a o      | False    |
| cac4966fa648df678b9f59117d085b40d647ef19 | b o      | False    |
| e0ca0a9d0afe2d168ed09efe2f859f76bcfd109f | c o      | True     |
+------------------------------------------+----------+----------+ `,
	RunE: func(cmd *cobra.Command, args []string) error {

		sort.Sort(sort.Reverse(migrations))

		err := migrations.Slice(batch.Hash, batch.Steps, Down)
		if err != nil {
			return err
		}

		if err := migrations.Execute(Down, db, batch.Exclude); err != nil {
			return err
		}
		return nil
	},
}