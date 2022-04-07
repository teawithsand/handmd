package tsrender

import "fmt"

// util function, which creates export default component with given tsx tags.
func TSExportDefaultComponent(tags string) string {
	return fmt.Sprintf(
		`export default () => {
return (
%s
);	
}`,
		tags,
	)
}
