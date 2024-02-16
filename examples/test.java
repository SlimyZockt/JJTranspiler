//Create a package
package examples;

//Import all the packages
import java.io.File;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

//Write a public parser class
class Main {
  String[] cars =  { "Volvo", "BMW", "Ford", "Mazda" };

  public static void main(String[] args) {
    try {
      // Access input data file and create document builder
      File inputDataFile = new File("InputData.txt");
      DocumentBuilderFactory dbldrFactory = DocumentBuilderFactory.newInstance();
      DocumentBuilder docBuilder = dbldrFactory.newDocumentBuilder();
      Document docmt = docBuilder.parse(inputDataFile);
      docmt.getDocumentElement().normalize();
      System.out.println("Name of the Root element:"
          + docmt.getDocumentElement().getNodeName());
      // Create node list
      NodeList ndList = docmt.getElementsByTagName("employee");
      System.out.println("*****************************");
      // Iterate through the node list and extract values
      for (int tempval = 0; tempval < ndList.getLength(); tempval++) {
        Node nd = ndList.item(tempval);
        System.out.println("\n Name of the current element :"
            + nd.getNodeName());
        if (nd.getNodeType() == Node.ELEMENT_NODE) {
          Element elemnt = (Element) nd;
          System.out.println("Employee ID : "
              + elemnt.getAttribute("empid"));
          System.out.println("Employee First Name: "
              + elemnt
                  .getElementsByTagName("firstname")
                  .item(0)
                  .getTextContent());
          System.out.println("Employee Last Name: "
              + elemnt
                  .getElementsByTagName("lastname")
                  .item(0)
                  .getTextContent());
          System.out.println("Employee Nick Name: "
              + elemnt
                  .getElementsByTagName("nickname")
                  .item(0)
                  .getTextContent());
          System.out.println("Employee Salary: "
              + elemnt
                  .getElementsByTagName("salary")
                  .item(0)
                  .getTextContent());
        }
      }
    } catch (Exception e) {
      // Catch and print exception - if any
      e.printStackTrace();
    }
 
  }
}


public class test {

  void test(){

  }

}