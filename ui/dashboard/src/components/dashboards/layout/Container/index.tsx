import Children from "../Children";
import ContainerTitle from "../../titles/ContainerTitle";
import Grid from "../Grid";
import { ContainerDefinition } from "@powerpipe/types";
import { ContainerProvider, useContainer } from "@powerpipe/hooks/useContainer";
import { registerComponent } from "../../index";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

type ContainerProps = {
  layoutDefinition?: ContainerDefinition;
  definition?: ContainerDefinition;
};

const Container = ({ definition }) => {
  const { showTitle } = useContainer();
  return (
    <Grid name={definition.name} width={definition.width}>
      {showTitle && <ContainerTitle title={definition.title} />}
      <Children
        childPanels={definition?.children || []}
        parentType="container"
      />
    </Grid>
  );
};

const ContainerWrapper = (props: ContainerProps) => {
  const { panelsMap } = useDashboardState();

  if (!props.definition && !props.layoutDefinition) {
    return null;
  }

  const panelDefinition = props.definition
    ? props.definition
    : props.layoutDefinition && panelsMap[props.layoutDefinition.name]
      ? panelsMap[props.layoutDefinition.name]
      : props.layoutDefinition;

  if (!panelDefinition) {
    return null;
  }

  return (
    <ContainerProvider>
      <Container
        definition={{
          ...panelDefinition,
          children: props.definition
            ? props.definition.children
            : props.layoutDefinition
              ? props.layoutDefinition.children
              : [],
        }}
      />
    </ContainerProvider>
  );
};

registerComponent("container", ContainerWrapper);

export default Container;
