package main

func (p *plugin) build() (undo bool, err error) {
	switch p.Type {
	case typeDeployment:
		undo, err = p.deployment()
	case typeStateful:
		undo, err = p.stateful()
	}

	return
}
